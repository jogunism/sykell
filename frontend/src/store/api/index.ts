import http from "@lib/http";
// constants
// import type { University, Contactus } from "@/types/constants";

export const crawlUrl = async (url: string) => {
  try {
    const response = await http.post("/crawl", { url: url });
    return response.data;
  } catch (error) {
    console.error("Error sending the url:", error);
    throw error;
  }
};
