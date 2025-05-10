"use client"

import SearchBar from "@/components/searchbar/searchbar"
import RecipeTree from "@/components/recipetree/tree"

const recipeData = {
  nodes: ["Acid Rain", "Rain", "Rain", "Smoke", "Smog"],
  dependencies: [
    { node: 0, dependency1: 1, dependency2: 3 },
    { node: 0, dependency1: 2, dependency2: 4 }
  ]
};

export default function Page() {
  return (
    <div className="bg-hero-pattern bg-center bg-black bg-blend-darken bg-opacity-70 flex min-h-screen flex-col items-center justify-start pt-10 gap-8">
      <div className="flex flex-col items-center h-1/4">
        <h1 className="text-6xl font-bold text-white mb-4">
          Little Alchemy 2 Recipe Finder
        </h1>
      </div>
      <SearchBar />
      <RecipeTree recipeJson={recipeData} />
    </div>
  )
}