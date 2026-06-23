export function InstructionSteps({ steps }: { steps: string[] }) {
  return <ol className="steps">{steps.map((step, index) => <li key={step}><span>{String(index + 1).padStart(2, "0")}</span><p>{step}</p></li>)}</ol>;
}
