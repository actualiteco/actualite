import { useEffect, useRef } from 'react';

interface UseInfiniteScrollOptions {
    onReach: () => void;
    enabled: boolean;
    rootMargin?: string;
}

export const useInfiniteScroll = ({
    onReach,
    enabled,
    rootMargin = '200px',
}: UseInfiniteScrollOptions) => {
    const sentinelRef = useRef<HTMLDivElement | null>(null);
    const callbackRef = useRef(onReach);

    useEffect(() => {
        callbackRef.current = onReach;
    });

    useEffect(() => {
        const node = sentinelRef.current;
        if (!enabled || !node) {
            return undefined;
        }

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries.some((entry) => entry.isIntersecting)) {
                    callbackRef.current();
                }
            },
            { rootMargin },
        );

        observer.observe(node);
        return () => observer.disconnect();
    }, [enabled, rootMargin]);

    return sentinelRef;
};
