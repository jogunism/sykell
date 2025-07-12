import React from "react";
import { create } from "zustand";
// APIs
import { crawlList, deleteCrawlItem, crawlUrl } from "./api";
// Utils
import { errorHandler } from "@lib/errorHandler";
import { toast } from "@lib/toast";

// constants
import type { CrawlItem, HeadingChartItem, LinkChartItem } from "@/constants";

/**
 * Main Store
 */
interface MainStore {
  pending: boolean;
  isSuccess: undefined | boolean;

  crawlItemList: CrawlItem[];
  currPage: number;
  pageSize: number;
  totalCount: number;
  fetchCrawlList: () => Promise<void>;

  checkedIds: number[];
  setCheckedIds: (ids: number[]) => void;
  clickCheckbox: (itemId: number, isChecked: boolean) => void;
  isAllChecked: boolean;

  queryString: string;
  sorting: { [key: string]: boolean | undefined };
  toggleSorting: (type: string) => Promise<void>;

  showDeleteButton: boolean;
  deleteCheckedItems: () => Promise<void>;

  reAnalyticsUrls: () => Promise<void>;

  setCurrPage: (currPage: number) => Promise<void>;
  setQueryString: (query: string) => void;

  currentItem: CrawlItem | null;
  setCurrentItem: (id?: number) => void;
  isModalOpen: boolean;
  setModalOpen: (isOpen: boolean) => void;

  headingChartData: HeadingChartItem[];
  linkChartData: LinkChartItem[];

  crawl: (url: string) => Promise<void>;
  reset: () => void;
}

const useMainStore = create<MainStore>((set, get) => ({
  /*************************************************************
   * Init State values
   *************************************************************/

  pending: false,
  isSuccess: undefined,

  crawlItemList: [],
  checkedIds: [],
  isAllChecked: false,

  sorting: {
    // page_title: true,
    // error: true,
    created_at: false,
  },

  currPage: 1,
  pageSize: 10,
  totalCount: 0,

  queryString: "",
  showDeleteButton: false,

  currentItem: null,
  isModalOpen: false,

  headingChartData: [],
  linkChartData: [],

  /*************************************************************
   * State handlers
   *************************************************************/

  /** HOME **/

  fetchCrawlList: async () => {
    set({ pending: true });

    try {
      const response = await crawlList(
        get().currPage,
        get().pageSize,
        get().queryString,
        JSON.stringify(get().sorting)
      );

      set({
        crawlItemList: response.list ?? [],
        totalCount: response.totalCount,
      });
    } catch (error) {
      const err = errorHandler(error);
      toast.error(err.message);
    }

    set({ pending: false });
  },

  setCheckedIds: (ids: number[]) => {
    set({ checkedIds: ids ?? [] });
    set({ showDeleteButton: ids?.length > 0 });
    set({
      isAllChecked:
        ids.length === get().crawlItemList.length &&
        get().crawlItemList.length > 0,
    });
  },

  clickCheckbox: (itemId: number, isChecked: boolean) => {
    let newCheckedIds = [...get().checkedIds];
    if (isChecked) {
      if (!newCheckedIds.includes(itemId)) {
        newCheckedIds.push(itemId);
      }
    } else {
      newCheckedIds = newCheckedIds.filter((id) => id !== itemId);
    }

    get().setCheckedIds(newCheckedIds);
  },

  toggleSorting: async (type: string) => {
    const currentSorting = get().sorting;
    const newSorting: { [key: string]: boolean | undefined } = {
      [type]:
        currentSorting[type] !== undefined ? !currentSorting[type] : false,
    };
    set({ sorting: newSorting });
    // console.log(newSorting);
    await get().fetchCrawlList();
  },

  setCurrPage: async (currPage: number) => {
    set({ currPage: currPage });

    await get().fetchCrawlList();
  },

  setQueryString: (query: string) => {
    set({ queryString: query });
  },

  deleteCheckedItems: async () => {
    set({ pending: true });
    try {
      const response = await deleteCrawlItem(get().checkedIds);

      // Re-fetch the list to ensure data consistency
      await get().fetchCrawlList();

      set({ isAllChecked: false });

      toast.success(response.message);
    } catch (error) {
      const err = errorHandler(error);
      toast.error(err.message);
    }

    set({ pending: false });
  },

  reAnalyticsUrls: async () => {
    set({ pending: true });

    const itemMap = new Map(get().crawlItemList.map((item) => [item.id, item]));

    try {
      const _promises = get()
        .checkedIds.map((id) => {
          const item = itemMap.get(id);
          if (item?.url) {
            return get().crawl(item.url);
          }
          console.warn(`NO url for ID ${id}. Skipping re-analysis.`);
          return null;
        })
        .filter((p): p is Promise<void> => p !== null); // Type guard to filter out nulls

      if (_promises.length > 0) {
        toast.info(`${_promises.length} items submitted for re-analysis.`);
        await Promise.all(_promises);
      }

      // Refresh the list to show updated data and clear selections
      await get().fetchCrawlList();
      get().setCheckedIds([]);
    } catch (error) {
      const err = errorHandler(error);
      toast.error(
        `An error occurred during the re-analysis batch process: ${err.message}`
      );
    } finally {
      set({ pending: false });
    }
  },

  setCurrentItem: (id?: number) => {
    const currItem = id
      ? get().crawlItemList.find((item: CrawlItem) => item.id === id) || null
      : null;

    set({ currentItem: currItem });
    if (currItem) {
      set({ isModalOpen: true });
      set({
        headingChartData: currItem.headingCounts
          ? Object.entries(currItem.headingCounts).map(([name, count]) => ({
              name: name.toUpperCase(), // e.g., H1, H2
              count: count,
            }))
          : [],
      });
      set({
        linkChartData: [
          {
            name: "Internal Links",
            value: currItem.internalLinkCount || 0,
          },
          {
            name: "External Links",
            value: currItem.externalLinkCount || 0,
          },
          {
            name: "Inaccessible Links",
            value: currItem.inaccessibleLinkCount || 0,
          },
        ].filter((entry) => entry.value > 0),
      });
    }
  },

  setModalOpen: (isOpen: boolean) => {
    set({ isModalOpen: isOpen });
  },

  /** ANALYTICS **/

  crawl: async (url: string) => {
    set({ pending: true });

    try {
      const response = await crawlUrl(url);

      const htmlContent = React.createElement("div", null, [
        response.message,
        React.createElement("br"),
        React.createElement("a", { href: `/?c=${response.id}` }, "See details"),
      ]);
      toast.success(htmlContent, { autoClose: false });

      set({ isSuccess: true });
      //
    } catch (error) {
      const err = errorHandler(error);
      toast.error(err.message);

      set({ isSuccess: false });
    }

    set({ pending: false });
  },

  reset: () => {
    set({ isSuccess: undefined });
  },
}));

export default useMainStore;
