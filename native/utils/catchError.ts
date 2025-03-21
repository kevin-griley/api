type Input<T> = Promise<T> | (() => T) | (() => Promise<T>);
type Output<T> = Promise<[undefined, T] | [unknown, undefined]>;

export async function catchError<T>(fn: Input<T>): Output<T> {
    try {
        const data = await (fn instanceof Promise
            ? fn
            : typeof fn === 'function'
              ? fn()
              : Promise.resolve(fn));
        return [undefined, data];
    } catch (error) {
        return [error, undefined];
    }
}
