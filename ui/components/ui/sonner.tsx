"use client"

import { Toaster as Sonner, type ToasterProps } from "sonner"

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      theme="light"
      className="toaster group"
      style={
        {
          "--normal-bg": "#000000",
          "--normal-text": "#ffffff",
          "--normal-border": "#000000",
          "--border-radius": "0.5rem",
        } as React.CSSProperties
      }
      toastOptions={{
        classNames: {
          toast: "!bg-black !text-white !border-black font-sans",
          description: "!text-neutral-300",
          actionButton: "!bg-white !text-black",
          cancelButton: "!bg-neutral-800 !text-white",
        },
      }}
      {...props}
    />
  )
}

export { Toaster }
