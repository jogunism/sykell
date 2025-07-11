import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

export default function ToastProvider() {
  return (
    <ToastContainer
      position="top-right"
      autoClose={5000}
      hideProgressBar={true}
      newestOnTop={false}
      closeOnClick
      rtl={false} // RTL (Right To Left) 모드 여부
      pauseOnFocusLoss
      draggable
      pauseOnHover
      theme="colored" // light, dark, colored
      // limit={3}
    />
  );
}
