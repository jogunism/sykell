import { create } from "zustand";
// APIs
import { crawlList, deleteCrawlItem, crawlUrl } from "./api";
// Utils
import { errorHandler } from "@lib/errorHandler";
import { toast } from "@lib/toast";
// constants
import type { CrawlItem } from "@/constants";

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
  showDeleteButton: boolean;

  setCurrPage: (currPage: number) => Promise<void>;
  setQueryString: (query: string) => void;

  deleteCheckedItems: () => Promise<void>;

  crawl: (url: string) => Promise<void>;
  reset: () => void;
}

const useMainStore = create<MainStore>((set, get) => ({
  // init getter values
  pending: false,
  isSuccess: undefined,

  crawlItemList: [],
  checkedIds: [],
  isAllChecked: false,

  currPage: 1,
  pageSize: 10,
  totalCount: 0,

  queryString: "",
  showDeleteButton: false,

  //
  fetchCrawlList: async () => {
    set({ pending: true });

    try {
      const response = await crawlList(
        get().currPage,
        get().pageSize,
        get().queryString
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

      toast.success(response.message);
    } catch (error) {
      const err = errorHandler(error);
      toast.error(err.message);
    }

    set({ pending: false });
  },

  crawl: async (url: string) => {
    set({ pending: true });

    try {
      const response = await crawlUrl(url);

      toast.success(response.message);
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
