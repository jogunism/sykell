import http from "@lib/http";
import { snakeToCamel } from "@/utils";
// constants
// import type { University, Contactus } from "@/types/constants";

export const crawlList = async (
  currPage: number,
  pageSize: number,
  queryString: string,
  sorting: string
) => {
  try {
    const response = await http.get(
      `/crawl/list?currPage=${currPage}&pageSize=${pageSize}&query=${queryString}&sorting=${sorting}`
    );
    return snakeToCamel(response?.data);
  } catch (error) {}
};

export const deleteCrawlItem = async (ids: number[]) => {
  try {
    const response = await http.delete("/crawl", { data: { ids: ids } });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const crawlUrl = async (url: string) => {
  try {
    const response = await http.post("/crawl", { url: url });
    return response.data;
  } catch (error) {
    console.error("Error sending the url:", error);
    throw error;
  }
};
