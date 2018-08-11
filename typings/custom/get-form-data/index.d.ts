interface getFormDataOptions {
  trim: boolean;
}

declare function getFormData (form: HTMLElement, options: getFormDataOptions): Array<string>;