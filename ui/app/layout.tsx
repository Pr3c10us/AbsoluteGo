import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import { GeistMono } from "geist/font/mono";
import { Providers } from "@/lib/providers";
import { Toaster } from "@/components/ui/sonner";
import Dock from "@/components/dock";
import UploadTracker from "@/components/upload-tracker";
import "./globals.css";

export const metadata: Metadata = {
  title: "AbsoluteGo",
  description: "Manage your comic book library, scripts, and videos",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${GeistSans.variable} ${GeistMono.variable} antialiased pb-24`}
      >
        <Providers>
          {children}
          <UploadTracker />
          <Dock />
          <Toaster />
        </Providers>
      </body>
    </html>
  );
}
