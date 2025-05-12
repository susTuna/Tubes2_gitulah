"use client"

import SearchBar from "@/components/searchbar/searchbar"
import RecipeTree from "@/components/recipetree/tree"

const recipeData = {
  nodes: ["Earth", "Air", "Dust", "Fire", "Gunpowder", "Fire", "Explosion", "Fire", "Fire", "Energy", "Atomic Bomb"],
  dependencies: [
    {node: 10, dependency1: 6, dependency2: 9},
    {node: 2, dependency1: 0, dependency2: 1},
    {node: 4, dependency1: 2, dependency2: 7},
    {node: 9, dependency1: 3, dependency2: 8},
    {node: 6, dependency1: 5, dependency2: 4}
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