import React from "react";
// store
import mainStore from "@store/mainStore";

/**
 * Component Table Searchbox
 */
const TableSearchbox: React.FC = () => {
  const {
    showDeleteButton,
    queryString,
    setQueryString,
    deleteCheckedItems,
    fetchCrawlList,
  } = mainStore();

  /*******************************************************
   * handlers
   */
  const handleSearchKeyUp = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      fetchCrawlList();
    }
  };

  const handleClearButtonClick = () => {
    setQueryString(""); // Clear search and re-fetch all items
    fetchCrawlList();
  };

  /*******************************************************
   * render
   */
  return (
    <div className="mb-3 flex justify-end items-center">
      <div className="relative flex items-center">
        <input
          type="text"
          placeholder="Search..."
          className="p-2 pl-10 pr-8 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-700 placeholder-gray-300"
          value={queryString}
          onChange={(e) => setQueryString(e.target.value)}
          onKeyUp={handleSearchKeyUp}
        />
        <svg
          className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          ></path>
        </svg>
        {queryString && (
          <button
            className="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
            onClick={handleClearButtonClick}
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </button>
        )}
      </div>
      {/* <button
        className="ml-2 w-[80px] h-10 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-md"
        onClick={() => handleSearch(queryString)}
      >
        Search
      </button> */}
      {showDeleteButton && (
        <button
          className="ml-2 w-[80px] h-10 bg-red-700 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-md"
          onClick={deleteCheckedItems}
        >
          Delete
        </button>
      )}
    </div>
  );
};

export default TableSearchbox;
