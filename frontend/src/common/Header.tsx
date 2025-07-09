import React from "react";

interface HeaderProps {
  toggleSidebar: () => void;
}

/**
 * Header
 */
const Header: React.FC<HeaderProps> = ({ toggleSidebar }) => {
  return (
    <header className="bg-white shadow-md p-4 flex justify-between items-center">
      {/* Hanberger Icon - mobile only */}
      <button
        className="md:hidden p-2 text-gray-700 focus:outline-none focus:bg-gray-100"
        onClick={toggleSidebar}
      >
        <svg
          className="w-6 h-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M4 6h16M4 12h16M4 18h16"
          ></path>
        </svg>
      </button>
      <div className="flex items-center ml-auto text-gray-800">Hyunwoo Cho</div>
    </header>
  );
};

export default Header;
