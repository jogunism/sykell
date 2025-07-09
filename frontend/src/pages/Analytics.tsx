import React, { useState, useRef, useEffect } from "react";
// Utils
import { FaSpinner } from "react-icons/fa";
// store
import mainStore from "@store/mainStore";

/**
 * Component Analytics
 */
const Analytics: React.FC = () => {
  const { crawl, pending, isSuccess } = mainStore();

  // states
  const [url, setUrl] = useState("");
  const [isValidUrl, setIsValidUrl] = useState(true);
  const [validationMessage, setValidationMessage] = useState<string | null>(
    null
  );

  const inputRef = useRef<HTMLInputElement>(null); // input box

  /*******************************************************
   * methods
   */
  const validateUrl = (value: string): boolean => {
    setIsValidUrl(true);
    setValidationMessage(null);

    if (value.trim() === "") {
      setIsValidUrl(false);
      setValidationMessage("Please input URL text.");
      return false;
    }

    if (!value.startsWith("http://") && !value.startsWith("https://")) {
      setIsValidUrl(false);
      setValidationMessage("Please input a valid URL text.");
      return false;
    }
    return true;
  };

  const handleSendClick = () => {
    if (validateUrl(url)) {
      crawl(url);
    }
  };

  const handleBlur = () => {
    validateUrl(url);
  };

  const resetInputValue = () => {
    setUrl("");
  };

  /*******************************************************
   * lifecycle hooks
   */
  useEffect(() => {
    // Focus the input element when the component mounts
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  useEffect(() => {
    if (isSuccess !== undefined) {
      resetInputValue();
    }
  }, [isSuccess]);

  /*******************************************************
   * render
   */
  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-800">Analytics</h1>
      <p className="mt-2 text-lg text-gray-600">
        Enter the URL for page analysis
      </p>
      <div className="mt-4 flex items-center">
        <input
          type="text"
          value={url}
          disabled={pending}
          onChange={(e) => {
            setUrl(e.target.value);
            setIsValidUrl(true);
            setValidationMessage(null);
          }}
          onBlur={handleBlur} // Add onBlur event handler
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              handleSendClick();
            }
          }}
          placeholder="https://url"
          className={`border p-2 rounded-md w-full max-w-md placeholder-gray-300 ${
            !isValidUrl ? "border-red-500" : "border-gray-500"
          } ${
            pending
              ? "bg-gray-100 cursor-not-allowed text-gray-500"
              : "bg-white text-gray-800"
          }
          }`}
          ref={inputRef} // Attach the ref to the input element
        />
        <button
          onClick={handleSendClick}
          className="ml-2 bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-md w-22 h-10 flex items-center justify-center"
        >
          {pending ? (
            <FaSpinner className="animate-spin" size={18} />
          ) : (
            <span>Send</span>
          )}
        </button>
      </div>
      {validationMessage && (
        <p className="text-red-500 text-sm mt-2">{validationMessage}</p>
      )}
    </div>
  );
};

export default Analytics;
