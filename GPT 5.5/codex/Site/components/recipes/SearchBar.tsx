import { Input } from "@/components/ui/Input";

export function SearchBar({ label, placeholder, defaultValue }: { label: string; placeholder: string; defaultValue?: string }) {
  return <label className="search-field"><span className="sr-only">{label}</span><span aria-hidden="true">⌕</span><Input name="search" defaultValue={defaultValue} placeholder={placeholder} /></label>;
}
