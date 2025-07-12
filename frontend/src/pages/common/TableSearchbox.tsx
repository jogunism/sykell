import React from "react";
// UI Component
import Magnifyer from "@icons/Magnifyer";
import X from "@icons/X";
// store
import mainStore from "@store/mainStore";

/**
 * Component Table Searchbox
 */
const TableSearchbox: React.FC = () => {
  const { queryString, setQueryString, fetchCrawlList } = mainStore();

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
        <Magnifyer />
        {queryString && (
          <button
            className="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
            onClick={handleClearButtonClick}
          >
            <X />
          </button>
        )}
      </div>
    </div>
  );
};

export default TableSearchbox;
