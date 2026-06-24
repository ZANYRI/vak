export function IngredientList({ ingredients }: { ingredients: string[] }) {
  return <ul className="ingredients">{ingredients.map((ingredient) => <li key={ingredient}>{ingredient}</li>)}</ul>;
}
