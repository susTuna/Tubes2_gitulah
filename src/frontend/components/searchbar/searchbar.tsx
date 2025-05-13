"use client"

import type React from "react"
import { useState } from "react"
import { Filter, Search, Check } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Checkbox } from "@/components/ui/checkbox"

interface FilterOption {
  id: string
  label: string
  checked: boolean
}

interface SearchBarProps {
  onSearch: (requestBody: unknown) => void
}

interface Element {
  id: number
  name: string
  tier: number
  image_url: string
}

export default function SearchBar({ onSearch }: SearchBarProps) {
  const [searchQuery, setSearchQuery] = useState("")
  const [isFilterOpen, setIsFilterOpen] = useState(false)
  const [multiCount, setMultiCount] = useState("1")
  const [delay, setDelay] = useState("1")
  const [filterOptions, setFilterOptions] = useState<FilterOption[]>([
    { id: "bfs", label: "BFS", checked: false },
    { id: "dfs", label: "DFS", checked: true  },
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
      }),
    )
  }

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault() // Prevent default form submission behavior
    console.log("Search button clicked")
    const activeFilters = filterOptions.filter((option) => option.checked).map((option) => option.label)

    console.log("Searching for:", searchQuery)
    console.log("With filters:", activeFilters)

    const selectedMethod = filterOptions.find((f) => (f.id === "bfs" || f.id === "dfs") && f.checked)?.id || "dfs"
    const threading = filterOptions.find((f) => f.id === "multi" && f.checked) ? "multi" : "single"

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_PUBLIC_API_URL}/elements/${searchQuery}?type=name`,
      )
      if (!response.ok) {
        throw new Error("Failed to fetch data from server")
      }

      const data = await response.json()
      console.log("Data received:", data)

      const element = data.find((item: Element) => item.name.toLowerCase() === searchQuery.toLowerCase())

      if (element) {
        console.log("Found element ID:", element.id)
      } else {
        console.log("Element not found")
      }

      const requestBody = {
        element: element.id,
        method: selectedMethod,
        count: filterOptions.find((f) => f.id === "multi" && f.checked) ? Number.parseInt(multiCount) : 1,
        delay: filterOptions.find((f) => f.id === "multi" && f.checked) ? Number.parseInt(delay) : 1,
        threading: threading,
      }

      onSearch(requestBody)
    } catch (error) {
      console.error("Error during search:", error)
    }
  }

  const isMultiSelected = filterOptions.find((f) => f.id === "multi")?.checked || false

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

      <Popover open={isFilterOpen} onOpenChange={setIsFilterOpen}>
        <PopoverTrigger asChild>
          <Button variant="outline" size="icon" className="text-black">
            <Filter className="h-4 w-4" />
            <span className="sr-only">Filter options</span>
          </Button>
        </PopoverTrigger>
        <PopoverContent align="end" className="w-56">
          <div className="space-y-4">
            <h4 className="font-medium">Method</h4>
            <div className="space-y-2">
              {filterOptions.map((option) => (
                <div key={option.id} className="flex items-center space-x-2">
                  <Checkbox
                    id={option.id}
                    checked={option.checked}
                    onCheckedChange={() => handleFilterChange(option.id)}
                  />
                  <Label htmlFor={option.id}>{option.label}</Label>
                </div>
              ))}
            </div>

            {isMultiSelected && (
              <div className="space-y-2">
                <Label htmlFor="count">Count:</Label>
                <Input
                  id="count"
                  type="number"
                  value={multiCount}
                  onChange={(e) => setMultiCount(e.target.value)}
                  min="1"
                />
                <Label htmlFor="delay">Delay:</Label>
                <Input
                  id="delay"
                  type="number"
                  value={delay}
                  onChange={(e) => setDelay(e.target.value)}
                  min="1"
                />
              </div>
            )}

            <div className="flex justify-end">
              <Button
                type="button"
                onClick={() => setIsFilterOpen(false)}
                className="flex items-center gap-2"
                size="sm"
              >
                <Check size={16} />
                Apply
              </Button>
            </div>
          </div>
        </PopoverContent>
      </Popover>

      <Button type="submit">Search</Button>
    </form>
  )
}
