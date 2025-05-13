"use client"

import { useState } from "react"
import Details from "@/components/details/details"
import SearchBar from "@/components/searchbar/searchbar"
import RecipeTree from "@/components/recipetree/tree"
import { fetchElementInfo } from "@/util/parser/parser"
import { fetchFromBackend } from "@/pages/api/proxy/[...path]"

export default function Page() {
  const [recipeData, setRecipeData] = useState(null)

  const handleSearch = async (requestBody: unknown) => {
    try {
      const response = await fetchFromBackend(`/fullrecipe/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(requestBody)
      })
  
      if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`)
  
      const data = await response.json()
      const searchId = data.search_id
  
      pollRecipeData(searchId)
    } catch (error) {
      console.error("Fetch error:", error)
    }
  }
  
  const pollRecipeData = async (searchId: number) => {
    const interval = setInterval(async () => {
      try {
        const response = await fetchFromBackend(`/fullrecipe/${searchId}`)
        if (!response.ok) throw new Error("Failed to fetch recipe data")
  
        const data = await response.json()
        
        const uniqueIds = Array.from(new Set<number>(data.nodes));
        const nameMap = await fetchElementInfo(uniqueIds);

        setRecipeData({ ...data, nameMap });
  
        if (data.finsihed) {
          clearInterval(interval)
        }
      } catch (error) {
        console.error("Polling error:", error)
        clearInterval(interval)
      }
    }, 250) // poll every 250ms
  }
  

  return (
    <div className="bg-hero-pattern bg-center bg-black bg-blend-darken bg-opacity-70 flex min-h-screen flex-col items-center justify-start pt-8 gap-8">
      <div className="flex flex-col items-center h-1/4">
        <h1 className="text-6xl font-bold text-white mb-4">
          Little Alchemy 2 Recipe Finder
        </h1>
      </div>
      <SearchBar onSearch={handleSearch} />
      {recipeData && <Details searchResults={recipeData} />}
      {recipeData && <RecipeTree recipeJson={recipeData} />}
    </div>
  )
}
