import React from "react";
// UI Component
import Modal from "./common/Modal";
// store
import mainStore from "@store/mainStore";
// utils
import { formatDate } from "@/utils";
// Recharts
import {
  ResponsiveContainer,
  BarChart,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  Bar,
  PieChart,
  Pie,
  Cell,
} from "recharts";

/**
 * Component Home detail modal
 */
const HomeDetail: React.FC = () => {
  const {
    currentItem,
    isModalOpen,
    setModalOpen,
    headingChartData,
    linkChartData,
  } = mainStore();

  const COLORS = ["#0088FE", "#00C49F", "#FFBB28"]; // Colors for pie chart segments

  /*******************************************************
   * handlers
   */
  const handleModalClose = () => {
    setModalOpen(false);
  };

  /*******************************************************
   * render
   */
  return (
    <Modal
      size={"lg"}
      isOpen={isModalOpen}
      onClose={handleModalClose}
      title={currentItem?.pageTitle}
      titleColor={currentItem?.error ? "text-red-800" : "text-gray-800"}
    >
      <div className="p-6 text-gray-700">
        <div className="p-5 mb-6 border border-gray-300 rounded-2xl lg:p-6 md:min-h-[200px]">
          <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div className="w-full">
              <h4
                className={`text-lg font-semibold ${
                  currentItem?.error ? "text-red-800" : "text-gray-800"
                } mb-6`}
              >
                Detail
              </h4>
              {/* detail */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-7">
                {currentItem?.error && (
                  <div className="col-span-full">
                    <p className="mb-2 text-xs leading-normal text-red-800 font-bold">
                      ERROR
                    </p>
                    <p className="text-sm font-medium text-red-800 ">
                      {currentItem?.error}
                    </p>
                  </div>
                )}
                <div>
                  <p className="mb-2 text-xs leading-normal text-gray-500 ">
                    URL
                  </p>
                  <p className="text-sm font-medium text-blue-800 ">
                    <a href={currentItem?.url || undefined} target="_blank">
                      {currentItem?.url}
                    </a>
                  </p>
                </div>

                <div>
                  <p className="mb-2 text-xs leading-normal text-gray-500 ">
                    Html version
                  </p>
                  <p className="text-sm font-medium text-gray-800 ">
                    {currentItem?.htmlVersion || "-"}
                  </p>
                </div>

                <div>
                  <p className="mb-2 text-xs leading-normal text-gray-500 ">
                    Has login form
                  </p>
                  <p className="text-sm font-medium text-gray-800 ">
                    {currentItem?.hasLoginForm ? "YES" : "NO"}
                  </p>
                </div>

                <div>
                  <p className="mb-2 text-xs leading-normal text-gray-500 ">
                    Date
                  </p>
                  <p className="text-sm font-medium text-gray-800 ">
                    {formatDate(currentItem?.createdAt) ?? "-"}
                  </p>
                </div>
              </div>
              {/* ./detail */}
            </div>
          </div>
        </div>

        <div className="p-5 mb-6 border border-gray-300 rounded-2xl lg:p-6 md:min-h-[200px]">
          <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div className="w-full">
              {/* chart */}
              <div className="flex flex-col md:flex-row md:gap-x-6">
                <div className="w-full md:w-1/2">
                  <h4 className="text-lg font-semibold text-gray-800 mb-4">
                    Heading Tag Counts
                  </h4>
                  {headingChartData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={headingChartData}>
                        <XAxis dataKey="name" />
                        <YAxis />
                        <Tooltip
                          contentStyle={{
                            backgroundColor: "transparent",
                            border: "none",
                          }}
                          wrapperStyle={{
                            outline: "none",
                            backgroundColor: "transparent",
                          }}
                        />
                        <Legend />
                        <Bar dataKey="count" fill="#8884d8" barSize={40} />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <p className="text-gray-500 text-center pt-10">
                      No heading data available.
                    </p>
                  )}
                </div>

                <div className="w-full md:w-1/2 mt-8 md:mt-0">
                  <h4 className="text-lg font-semibold text-gray-800 mb-4">
                    Link Counts
                  </h4>
                  {linkChartData.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <PieChart>
                        <Pie
                          data={linkChartData}
                          cx="50%"
                          cy="50%"
                          innerRadius={0}
                          outerRadius={80}
                          fill="#8884d8"
                          paddingAngle={0}
                          dataKey="value"
                        >
                          {linkChartData.map((_, index) => (
                            <Cell
                              key={`cell-${index}`}
                              fill={COLORS[index % COLORS.length]}
                            />
                          ))}
                        </Pie>
                        <Tooltip />
                        <Legend />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <p className="text-gray-500 text-center pt-10">
                      No link data available.
                    </p>
                  )}
                </div>
              </div>
              {/* ./chart */}
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default HomeDetail;
