import { create } from "zustand";
// APIs
import { crawlUrl } from "./api";
// Utils
import { errorHandler } from "@lib/errorHandler";
import { toast } from "@lib/toast";
// // constants
// import type { Contactus } from '@/types/constants';

interface MainStore {
  pending: boolean;
  isSuccess: undefined | boolean;
  errorMessage: string | null; // New state for error message
  crawl: (url: string) => Promise<void>;
  reset: () => void;
}

const useMainStore = create<MainStore>((set) => ({
  pending: false,
  isSuccess: undefined,
  errorMessage: null, // Initialize error message

  crawl: async (url: string) => {
    set({ pending: true, errorMessage: null }); // Reset error message on new crawl

    try {
      const response = await crawlUrl(url);

      toast.success(response.message);
      set({ isSuccess: true });
      //
    } catch (error) {
      const err = errorHandler(error);

      toast.error(err.message);

      set({ isSuccess: false });
    }

    set({ pending: false });
  },

  reset: () => {
    set({ isSuccess: undefined });
  },
}));

export default useMainStore;
