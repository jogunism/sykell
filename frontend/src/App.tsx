import { useState } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
// Components
import Header from "@/pages/common/Header";
import LeftMenu from "@/pages/common/LeftMenu";
import Home from "@/pages/Home";
import Analytics from "@/pages/Analytics";

/**
 * App
 */
function App() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <BrowserRouter>
      <div className="flex flex-col md:flex-row h-screen bg-gray-100">
        <LeftMenu
          isSidebarOpen={isSidebarOpen}
          setIsSidebarOpen={setIsSidebarOpen}
        />

        <div className="flex-1 flex flex-col">
          {/* Header */}
          <Header toggleSidebar={toggleSidebar} />

          {/* Main Content */}
          <main className="flex-1 p-6 bg-gray-100 overflow-y-auto">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/analytics" element={<Analytics />} />
            </Routes>
          </main>
        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
