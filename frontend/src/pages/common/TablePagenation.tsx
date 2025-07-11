import React from "react";
// store
import mainStore from "@store/mainStore";

/**
 * Component Table Pagenation
 */
const TablePagenation: React.FC = () => {
  const { currPage, setCurrPage, pageSize, totalCount } = mainStore();

  const totalPages = Math.ceil(totalCount / pageSize);
  const pageNumbers = Array.from({ length: totalPages }, (_, i) => i + 1);

  return (
    <div className="px-5 py-5 bg-white border-t flex flex-col xs:flex-row items-center xs:justify-between">
      {/* <span className="text-xs xs:text-sm text-gray-900">
        Showing {(currPage - 1) * pageSize + 1} to{" "}
        {Math.min(currPage * pageSize, totalCount)} of {totalCount} Entries
      </span> */}
      <div className="inline-flex mt-2 xs:mt-0 rounded-md shadow-sm -space-x-px">
        <button
          onClick={() => setCurrPage(currPage - 1)}
          disabled={currPage === 1}
          className="relative inline-flex items-center px-4 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
        >
          Prev
        </button>
        {pageNumbers.map((page) => (
          <button
            key={page}
            onClick={() => setCurrPage(page)}
            className={`relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium ${
              currPage === page
                ? "bg-blue-500 text-white"
                : "bg-white text-gray-700 hover:bg-gray-50"
            }`}
          >
            {page}
          </button>
        ))}
        <button
          onClick={() => setCurrPage(currPage + 1)}
          disabled={currPage === totalPages}
          className="relative inline-flex items-center px-4 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default TablePagenation;
