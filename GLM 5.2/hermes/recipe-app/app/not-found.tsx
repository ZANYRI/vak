import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/Button";

export default function GlobalNotFound() {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center px-4 text-center">
      <p className="font-serif text-7xl font-semibold text-primary">404</p>
      <h1 className="mt-4 font-serif text-2xl font-semibold text-foreground">
        Page not found
      </h1>
      <p className="mt-2 text-muted">
        The page you are looking for does not exist.
      </p>
      <div className="mt-6">
        <Link href="/">
          <Button>Go home</Button>
        </Link>
      </div>
    </div>
  );
}
