import React from "react";
import { Link } from "react-router-dom";

interface LeftMenuProps {
  isSidebarOpen: boolean;
  setIsSidebarOpen: (isOpen: boolean) => void;
}

/**
 * Leftmenus
 */
const LeftMenu: React.FC<LeftMenuProps> = ({
  isSidebarOpen,
  setIsSidebarOpen,
}) => {
  return (
    <>
      {/* Sidebar over rayout for Mobile */}
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-translate z-40 md:hidden"
          onClick={() => setIsSidebarOpen(false)}
        ></div>
      )}

      <aside
        className={`fixed inset-y-0 left-0 w-64 bg-white shadow-lg flex flex-col z-50
          transform ${isSidebarOpen ? "translate-x-0" : "-translate-x-full"}
          transition-transform duration-300 ease-in-out
          md:relative md:translate-x-0 md:flex`}
      >
        <div className="p-6 border-b border-gray-200">
          <Link
            to="/"
            className="text-2xl font-bold text-gray-800"
            onClick={() => setIsSidebarOpen(false)}
          >
            Sykell
          </Link>
        </div>
        <nav className="flex-1 p-4">
          <ul>
            <li className="mb-2">
              <Link
                to="/"
                className="flex items-center p-2 text-gray-700 hover:bg-gray-100 rounded-md"
                onClick={() => setIsSidebarOpen(false)}
              >
                <span className="mr-3">📊</span> Home
              </Link>
            </li>
            <li className="mb-2">
              <Link
                to="/analytics"
                className="flex items-center p-2 text-gray-700 hover:bg-gray-100 rounded-md"
                onClick={() => setIsSidebarOpen(false)}
              >
                <span className="mr-3">📈</span> Analytics
              </Link>
            </li>
          </ul>
        </nav>
      </aside>
    </>
  );
};

export default LeftMenu;
