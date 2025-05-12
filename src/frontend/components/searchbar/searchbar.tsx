import type React from "react"

import { useState } from "react"
import { Filter, Search } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Input } from "@/components/ui/input"

interface FilterOption {
  id: string
  label: string
  checked: boolean
}

export default function SearchBar() {
  const [searchQuery, setSearchQuery] = useState("")
  const [filterOptions, setFilterOptions] = useState<FilterOption[]>([
    { id: "bfs", label: "BFS", checked: true },
    { id: "dfs", label: "DFS", checked: false },
    { id: "multi", label: "Multiple", checked: false },
  ])

  const handleFilterChange = (id: string) => {
    setFilterOptions((prevFilterOptions) =>
      prevFilterOptions.map((option) => {
        if (option.id === id) {
          // Toggle the checked state of the clicked filter option
          return { ...option, checked: !option.checked }
        }
        
        // Only uncheck "bfs" or "dfs" when one of them is clicked
        if ((option.id === "bfs" || option.id === "dfs") && option.checked && (id === "bfs" || id === "dfs")) {
          return { ...option, checked: false }
        }
        
        return option
      })
    )
  }
  
  

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault() // Prevent default form submission behavior
    console.log("Search button clicked")
    const activeFilters = filterOptions.filter((option) => option.checked).map((option) => option.label)

    console.log("Searching for:", searchQuery)
    console.log("With filters:", activeFilters)
    // Here you would typically call your search function with these parameters
  }

  return (
    <form onSubmit={handleSearch} className="flex w-full max-w-lg items-center space-x-2 text-white">
      <div className="relative flex-1">
        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          type="search"
          placeholder="Search..."
          className="pl-8"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="icon" className="text-black">
            <Filter className="h-4 w-4" />
            <span className="sr-only">Filter options</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>Method</DropdownMenuLabel>
          <DropdownMenuSeparator />
          {filterOptions.map((option) => (
            <DropdownMenuCheckboxItem
              key={option.id}
              checked={option.checked}
              onCheckedChange={() => handleFilterChange(option.id)}
            >
              {option.label}
            </DropdownMenuCheckboxItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>

      <Button type="submit">Search</Button>
    </form>
  )
}