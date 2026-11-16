// Calendar dates ("2019-04-02") travel as plain strings between the API and
// the UI. These helpers convert in *local* time. new Date("2019-04-02") and
// date.toISOString() both use UTC, which can shift a date by a day depending
// on the user's timezone.

/** "2019-04-02" -> Date at local midnight on 2 April 2019. */
export function parseIsoDate(value: string): Date {
    const [year, month, day] = value.split('-').map(Number);
    return new Date(year!, month! - 1, day!);
}

/** Date -> "2019-04-02", using the local calendar day. */
export function toIsoDate(date: Date): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}
