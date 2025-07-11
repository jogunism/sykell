// constants

export interface CrawlItem {
  id: number;
  htmlVersion?: string;
  pageTitle?: string;
  headingCounts?: HeadingCounts;
  internalLinkCount?: number;
  externalLinkCount?: number;
  inaccessibleLinkCount?: number;
  hasLoginForm?: boolean;
  url?: string | null;
  error?: string;
  createdAt?: string;
}

export interface HeadingCounts {
  h1: number;
  h2: number;
  h3: number;
  h4: number;
  h5: number;
  h6: number;
}

export interface HeadingChartItem {
  name: string;
  count: number;
}

export interface LinkChartItem {
  name: string;
  value: number;
}
