import i18n from 'sveltekit-i18n';

/** @type {import('sveltekit-i18n').Config} */
export const config = {
  loaders: [
    {
      locale: 'en',
      key: 'c',
      loader: async () => (await import('./en.json')).default,
    },
    {
      locale: 'de',
      key: 'c',
      loader: async () => (await import('./de.json')).default,
    },
    {
      locale: 'fr',
      key: 'c',
      loader: async () => (await import('./fr.json')).default,
    }
  ],
};

export const { t, loading, locales, locale, loadTranslations } = new i18n(config);