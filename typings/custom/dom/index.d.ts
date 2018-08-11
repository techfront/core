declare class Dom {
    static Ready(a: () => void): void;
    static Exists(a: string): boolean;
    static All(a: string): NodeListOf<Element>;
    static Nearest(a: Element | Node, b: string): NodeListOf<Element> | void;
    static First(a: string): Element;
    static Each(a: string, b: Function): void;
    static ForEach(a: Array<any>, b: Function): void;
    static Hide(a: string): void;
    static Show(a: string): void;
    static ShowHide(a: string): void;
    static On(a: string, b: string, c: EventListenerOrEventListenerObject): void;
    static Format(a: string): string;
    static Post(a: string, b: any, c: Function, d: Function): void;
    static Get(a: string, b: Function, c: any): void;
}