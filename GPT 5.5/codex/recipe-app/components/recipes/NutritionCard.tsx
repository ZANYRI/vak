export function NutritionCard({ title, labels, values }: { title: string; labels: [string, string, string, string]; values: [number, string, string, string] }) {
  return <section className="nutrition"><h2>{title}</h2><div>{labels.map((label, index) => <div key={label}><strong>{values[index]}</strong><span>{label}</span></div>)}</div></section>;
}
