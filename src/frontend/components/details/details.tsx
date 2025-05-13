'use client'

import { Card } from "@/components/ui/card"
import { RecipeJson } from "@/util/parser/parser"

export default function Details({searchResults} : {searchResults : RecipeJson}) {
    return ( 
        <Card className="flex w-full max-w-lg bg-transparent text-white items-center justify-center">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-20">
            <div className="flex flex-col items-center">
                <span className="text-gray-300 text-sm">Nodes searched: </span>
                <span className="font-semibold">{searchResults.nodes_searched} nodes</span>
            </div>
            <div className="flex flex-col items-center">
                <span className="text-gray-300 text-sm">Time taken: </span>
                <span className="font-semibold">{searchResults.time_taken} ms</span>
            </div>
            <div className="flex flex-col items-center">
                <span className="text-gray-300 text-sm">Recipes found: </span>
                <span className="font-semibold">{searchResults.recipes_found}</span>
            </div>
            </div>
        </Card>
    )
}
