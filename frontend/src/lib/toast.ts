import { ToastContainer, toast } from "react-toastify";

export interface ToastProps {
  hideProgressBar: boolean;
  autoClose: number;
  closeOnClick: boolean;
  pauseOnHover: boolean;
}

export { ToastContainer, toast };

// ref.
// https://fkhadra.github.io/react-toastify
