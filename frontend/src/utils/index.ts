// Utils

/**
 * Change snake_case to camelCase
 * @param obj
 * @returns obj
 */
export const snakeToCamel = (obj: any): any => {
  if (Array.isArray(obj)) {
    return obj.map((v) => snakeToCamel(v));
  } else if (obj !== null && typeof obj === "object") {
    return Object.keys(obj).reduce((acc, key) => {
      const camelKey = key.replace(/_([a-z])/g, (g) => g[1].toUpperCase());
      acc[camelKey] = snakeToCamel(obj[key]);
      return acc;
    }, {} as any);
  }
  return obj;
};

/**
 * Formats a date string to 'yyyy-MM-dd hh:mm' format.
 * @param dateString The date string in ISO 8601 format (e.g., '2025-07-10T20:01:16Z').
 * @returns Formatted date string or empty string if invalid.
 */
export const formatDate = (dateString?: string): string => {
  if (!dateString) return "";
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return ""; // Check for invalid date

  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}`;
};
