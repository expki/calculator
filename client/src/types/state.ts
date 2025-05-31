import type { State } from '@lib/schema/state';

export type LocalState = Partial<{
    Id: number,
    X: number,
    Y: number,
    CpuLoad: number, // render cpu utilization ratio (0.0 - 1.0)
    State: State,
}>;
