import React, { useEffect, useRef } from "react";
import X from "@icons/X";

interface ModalProps {
  isOpen: boolean;
  onClose?: () => void;
  children: React.ReactNode;
  title?: string;
  isParentBlur?: boolean;
  size?: "sm" | "md" | "lg" | "full";
  isScrollAllowed?: boolean;
  position?: "top" | "middle" | "bottom";
  titleColor?: string;
}

const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  children,
  title,
  isParentBlur = true,
  size = "sm",
  isScrollAllowed = false,
  position = "middle",
  titleColor = "text-gray-700",
}) => {
  const modalRef = useRef<HTMLDivElement>(null);

  const positionClasses = {
    top: "items-start",
    middle: "items-center",
    bottom: "items-end",
  };

  const backdropClasses = `fixed inset-0 z-50 flex justify-center ${
    positionClasses[position]
  } ${isParentBlur ? "backdrop-blur-sm backdrop-brightness-50" : ""}`;

  // size prop에 따른 max-w-* 클래스 매핑
  const sizeClasses = {
    sm: "w-full md:max-w-lg",
    md: "w-full md:max-w-3xl",
    lg: "w-full md:max-w-5xl",
    full: "w-screen h-screen rounded-none shadow-none",
  };

  /*******************************************************
   * lifecycle hooks
   */
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape" && onClose) {
        onClose();
      }
    };

    if (isOpen) {
      if (!isScrollAllowed) {
        document.body.style.overflow = "hidden";
      }
      document.addEventListener("keydown", handleEscape);
    } else {
      document.body.style.overflow = "";
      document.removeEventListener("keydown", handleEscape);
    }

    return () => {
      document.body.style.overflow = "";
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isOpen, onClose, isScrollAllowed]);

  /*******************************************************
   * render
   */
  if (!isOpen) return null;

  return (
    <div
      className={backdropClasses}
      onClick={(e) => {
        if (
          modalRef.current &&
          !modalRef.current.contains(e.target as Node) &&
          onClose
        ) {
          onClose();
        }
      }}
    >
      <div
        ref={modalRef}
        className={`relative bg-white rounded-lg shadow-xl ${sizeClasses[size]} flex flex-col max-h-[90vh]`}
      >
        <div className="flex-shrink-0">
          {onClose && (
            <button
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-700 text-2xl p-2 mt-2 mr-2"
              onClick={onClose}
            >
              <X />
            </button>
          )}
          {title && (
            <h1 className={`text-2xl font-bold px-6 py-4 ${titleColor}`}>
              {title}
            </h1>
          )}
        </div>
        <div className="flex-grow overflow-y-auto">{children}</div>
      </div>
    </div>
  );
};

export default Modal;
