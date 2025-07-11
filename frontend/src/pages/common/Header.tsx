import React from "react";
import Hamberger from "@icons/Hanberger";

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
        <Hamberger />
      </button>
      <div className="flex items-center ml-auto text-gray-800">Hyunwoo Cho</div>
    </header>
  );
};

export default Header;
