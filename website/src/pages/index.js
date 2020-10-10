import React from 'react';
import clsx from 'clsx';
import Layout from '@theme/Layout';
import useThemeContext from '@theme/hooks/useThemeContext';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';

const features = [
  {
    title: 'Easy to Use',
    imageUrl: 'img/command-window.svg',
    description: (
      <>
        No installation required.
        Does not require setting the PATH.
      </>
    ),
  },
  {
    title: 'No dependencies',
    imageUrl: 'img/unlink-symbol.svg',
    description: (
      <>
        Fisherman is distributed as a binary executable file.
        You don't need anything else to work.
      </>
    ),
  },
  {
    title: 'Flexible',
    imageUrl: 'img/configuration-gears.svg',
    description: (
      <>
        Simple but very flexible declarative configuration, that solves common issues in a couple of lines.
      </>
    ),
  },
];

const Feature = ({ imageUrl, title, description }) => {
  const imgUrl = useBaseUrl(imageUrl);
  return (
    <div className={clsx('col col--4 text--center', styles.feature)}>
      {imgUrl && (
        <div>
          <img className={styles.featureImage} src={imgUrl} alt={title} />
        </div>
      )}
      <h3>{title}</h3>
      <p className={styles.featureDescription}>{description}</p>
    </div>
  );
}

const Image = ({ image, darkImage, className }) => {
  const { isDarkTheme } = useThemeContext();
  return <img src={useBaseUrl(isDarkTheme ? darkImage : image)} className={className} />;
};

function Home() {
  const context = useDocusaurusContext();
  const { siteConfig = {} } = context;
  return (
    <Layout
      title={siteConfig.title}
      description="Description will go into a meta tag in <head />">
      <header className={clsx('hero hero--primary', styles.heroBanner)}>
        <div className="container">
          <div>
            <Image image="img/preview.png" darkImage="img/preview_dark.png" className={styles.heroBannerImage} />
          </div>
          <p className={clsx("hero__subtitle", styles.heroSubtitle)}>{siteConfig.tagline}</p>
          <div className={styles.buttons}>
            <Link
              className={clsx(
                'button button--outline button--lg',
                styles.getStarted,
              )}
              to={useBaseUrl('/docs/getting-started')}>
              Get Started
            </Link>
          </div>
        </div>
      </header>
      <main>
        {features && features.length > 0 && (
          <section className={styles.features}>
            <div className="container">
              <div className="row">
                {features.map((props, idx) => (
                  <Feature key={idx} {...props} />
                ))}
              </div>
            </div>
          </section>
        )}
      </main>
    </Layout>
  );
}

export default Home;
