import { NextResponse } from 'next/server';

// Always evaluated at request time so it reflects real liveness.
export const dynamic = 'force-dynamic';

export function GET() {
  return NextResponse.json(
    { status: 'ok', service: 'recipe-app' },
    { status: 200 },
  );
}
