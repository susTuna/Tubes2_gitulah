"use client"

import SearchBar from "@/components/searchbar/searchbar"

export default function Page() {
  return (
    <div className="bg-hero-pattern bg-center bg-black bg-blend-darken bg-opacity-70 flex min-h-screen flex-col items-center justify-start pt-[30vh] p-4 gap-8">
      <h1 className="text-6xl font-bold text-white mb-4">Gitulah</h1>
      <SearchBar />
    </div>
  )
}