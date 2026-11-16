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

/**
 * A pet's age from its "YYYY-MM-DD" date of birth, e.g. "3 years, 2 months",
 * "5 months" or "Under a month". Returns '' when the date isn't known.
 */
export function formatAge(dateOfBirth: string | null, today: Date = new Date()): string {
    if (!dateOfBirth) {
        return '';
    }
    const birthDate = parseIsoDate(dateOfBirth);

    let totalMonths = (today.getFullYear() - birthDate.getFullYear()) * 12
        + (today.getMonth() - birthDate.getMonth());
    // The birthday hasn't come round yet this month.
    if (today.getDate() < birthDate.getDate()) {
        totalMonths--;
    }
    if (totalMonths < 1) {
        return 'Under a month';
    }

    const years = Math.floor(totalMonths / 12);
    const months = totalMonths % 12;
    const plural = (count: number, word: string) => `${count} ${word}${count === 1 ? '' : 's'}`;

    if (years === 0) {
        return plural(months, 'month');
    }
    if (months === 0) {
        return plural(years, 'year');
    }
    return `${plural(years, 'year')}, ${plural(months, 'month')}`;
}

const relativeFormat = new Intl.RelativeTimeFormat(undefined, {numeric: 'auto'});

/**
 * How long ago a timestamp was, e.g. "3 hours ago", "yesterday", "5 days ago".
 * Older than a few weeks falls back to the date itself.
 */
export function formatTimeAgo(timestamp: string, now: Date = new Date()): string {
    const seconds = Math.round((new Date(timestamp).getTime() - now.getTime()) / 1000);
    const minutes = Math.round(seconds / 60);
    const hours = Math.round(minutes / 60);
    const days = Math.round(hours / 24);

    if (Math.abs(seconds) < 60) {
        return 'just now';
    }
    if (Math.abs(minutes) < 60) {
        return relativeFormat.format(minutes, 'minute');
    }
    if (Math.abs(hours) < 24) {
        return relativeFormat.format(hours, 'hour');
    }
    if (Math.abs(days) < 28) {
        return relativeFormat.format(days, 'day');
    }
    return new Date(timestamp).toLocaleDateString(undefined, {day: 'numeric', month: 'short', year: 'numeric'});
}
