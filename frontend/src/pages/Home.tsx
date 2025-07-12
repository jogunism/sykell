import React, { useEffect } from "react";
import { useSearchParams } from "react-router-dom";
// UI Component
import TablePagenation from "./common/TablePagenation";
import TableSearchbox from "./common/TableSearchbox";
import HomeDetail from "./HomeDetail";
import ArrowUp from "@icons/ArrowUp";
import ArrowDown from "@icons/ArrowDown";
// store
import mainStore from "@store/mainStore";
// constants
import type { CrawlItem } from "@/constants";
// utils
import { formatDate } from "@/utils";

/**
 * Component Home
 */
const Home: React.FC = () => {
  const {
    pending,
    crawlItemList,
    fetchCrawlList,

    showDeleteButton,
    deleteCheckedItems,
    reAnalyticsUrls,

    checkedIds,
    setCheckedIds,
    clickCheckbox,
    isAllChecked,

    sorting,
    toggleSorting,

    setCurrentItem,
    isModalOpen,
  } = mainStore();

  const [searchParams, setSearchParams] = useSearchParams();
  const currentId = searchParams.get("c");

  /*******************************************************
   * handlers
   */
  const handleCheckboxAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    const isChecked = e.target.checked;
    setCheckedIds(isChecked ? crawlItemList.map((item) => item.id) : []);
  };

  const handleCheckbox = (
    e: React.ChangeEvent<HTMLInputElement>,
    itemId: number
  ) => {
    clickCheckbox(itemId, e.target.checked);
  };

  const handleItemClick = (id: number) => {
    setCurrentItem(id);
  };

  /*******************************************************
   * lifecycle hooks
   */
  useEffect(() => {
    fetchCrawlList();
  }, []);

  useEffect(() => {
    if (currentId) {
      setTimeout(() => {
        setCurrentItem(parseInt(currentId));
      }, 200);
    }
  }, [currentId]);

  useEffect(() => {
    if (!isModalOpen) {
      // remove query parameter
      const newParams = new URLSearchParams(searchParams);
      newParams.delete("c");
      setSearchParams(newParams);
    }
  }, [isModalOpen]);

  /*******************************************************
   * render
   */
  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Home</h1>
      </div>

      <div className="mt-4 p-4 w-full bg-white shadow-md rounded-lg overflow-hidden ">
        <div className="flex items-center mb-3">
          <div className="flex items-center gap-2 pb-3">
            {showDeleteButton && (
              <>
                <button
                  className="w-[60px] h-10 bg-red-500 hover:bg-red-700 text-white text-sm font-bold py-2 px-2 rounded-md"
                  onClick={deleteCheckedItems}
                >
                  Delete
                </button>
                <button
                  className="w-[70px] h-10 bg-blue-500 hover:bg-blue-700 text-white text-sm font-bold py-2 px-1 rounded-md"
                  onClick={reAnalyticsUrls}
                >
                  Anaytics
                </button>
              </>
            )}
          </div>
          <div className="ml-auto">
            <TableSearchbox />
          </div>
        </div>

        <table className="min-w-full leading-normal table-fixed">
          <thead>
            <tr>
              <th className="px-2 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider w-8">
                <input
                  type="checkbox"
                  className="form-checkbox h-4 w-4 text-blue-600 transition duration-150 ease-in-out"
                  checked={isAllChecked}
                  onChange={handleCheckboxAll}
                />
              </th>
              <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider w-[60%] md:w-3/5">
                <div className="flex items-left">
                  <button
                    className="flex"
                    onClick={() => toggleSorting("page_title")}
                  >
                    <span>Page title</span>
                    {sorting.page_title === false ? <ArrowDown /> : <ArrowUp />}
                  </button>
                </div>
              </th>
              <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider hidden md:table-cell md:w-[10%]">
                Link
              </th>
              <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider w-[40%] md:w-[10%]">
                <div className="flex items-center justify-center">
                  <button
                    className="flex"
                    onClick={() => toggleSorting("error")}
                  >
                    <span>Status</span>
                    {sorting.error === false ? <ArrowDown /> : <ArrowUp />}
                  </button>
                </div>
              </th>
              <th className="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider hidden md:table-cell md:w-[20%]">
                <div className="flex items-center justify-center">
                  <button
                    className="flex"
                    onClick={() => toggleSorting("created_at")}
                  >
                    <span>Date</span>
                    {sorting.created_at === false ? <ArrowDown /> : <ArrowUp />}
                  </button>
                </div>
              </th>
            </tr>
          </thead>
          <tbody>
            {pending ? (
              <tr>
                <td
                  colSpan={5}
                  className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center text-gray-500"
                >
                  Loading...
                </td>
              </tr>
            ) : crawlItemList.length < 1 ? (
              <tr>
                <td
                  colSpan={5}
                  className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center text-gray-500"
                >
                  No Items
                </td>
              </tr>
            ) : (
              crawlItemList.map((item: CrawlItem) => {
                const isSuccess: boolean = item.error === "";
                return (
                  <tr key={item.id}>
                    <td className="px-2 py-5 border-b bg-white text-sm">
                      <input
                        type="checkbox"
                        className="form-checkbox h-4 w-4 text-blue-600 transition duration-150 ease-in-out"
                        checked={checkedIds.includes(item.id)}
                        onChange={(e) => handleCheckbox(e, item.id)}
                      />
                    </td>
                    <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                      <p
                        className={`text-gray-900 whitespace-no-wrap ${
                          !isSuccess ? "text-red-700" : ""
                        }`}
                      >
                        <span
                          onClick={() => handleItemClick(item.id)}
                          className="cursor-pointer"
                        >
                          {item.pageTitle}
                        </span>
                      </p>
                      {item.url && (
                        <button
                          onClick={() =>
                            window.open(
                              item.url || "",
                              "_blank",
                              "noopener noreferrer"
                            )
                          }
                          className={`mt-2 md:hidden inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white ${
                            isSuccess
                              ? "bg-blue-500 hover:bg-blue-700"
                              : "bg-red-400"
                          }`}
                        >
                          LINK
                        </button>
                      )}
                      <div className="mt-1 md:hidden text-xs text-gray-500">
                        {formatDate(item.createdAt)}
                      </div>
                    </td>
                    <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm hidden md:table-cell text-center">
                      {item.url && (
                        <button
                          onClick={() =>
                            window.open(
                              item.url || "",
                              "_blank",
                              "noopener noreferrer"
                            )
                          }
                          className={`inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded shadow-sm text-white ${
                            isSuccess
                              ? "bg-blue-500 hover:bg-blue-700"
                              : "bg-red-400"
                          }`}
                        >
                          LINK
                        </button>
                      )}
                    </td>
                    <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center">
                      <span
                        className={`relative inline-block px-3 py-1 font-semibold leading-tight ${
                          isSuccess ? "text-green-900" : "text-red-900"
                        }`}
                      >
                        <span
                          aria-hidden
                          className={`absolute inset-0 opacity-50 rounded-full ${
                            isSuccess ? "bg-green-400" : "bg-red-200"
                          }`}
                        />
                        <span className="relative text-xs">
                          {isSuccess ? "Success" : "Failed"}
                        </span>
                      </span>
                    </td>
                    <td className="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center hidden md:table-cell">
                      <p className="text-gray-900 whitespace-no-wrap">
                        {formatDate(item.createdAt)}
                      </p>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>

        {/* Pagination */}
        <TablePagenation />
      </div>

      <HomeDetail />
    </div>
  );
};

export default Home;
