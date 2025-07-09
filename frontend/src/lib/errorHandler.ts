import { AxiosError } from "axios";

interface HandledError {
  message: string;
  isAxiosError: boolean;
  statusCode?: number;
  responseBody?: any;
}

export function errorHandler(error: unknown): HandledError {
  let errorMessage: string;
  let isAxiosError = false;
  let statusCode: number | undefined;
  let responseBody: any | undefined;

  if (error instanceof AxiosError) {
    isAxiosError = true;
    statusCode = error.response?.status;
    responseBody = error.response?.data;

    errorMessage = error.response?.data?.message || error.message;
    console.error("Axios Error:", error.response?.data || error.message);
    console.error(error.stack);
  } else if (error instanceof Error) {
    errorMessage = error.message;
    console.error("Generic Error:", error.message);
    console.error(error.stack);
  } else {
    errorMessage = "An unknown error occurred.";
    console.error("Unknown Error Type:", error);
  }

  return {
    message: errorMessage,
    isAxiosError,
    statusCode,
    responseBody,
  };
}
