import React from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import clsx from 'clsx';
import Layout from '@theme/Layout';

function Hello() {
  return (
    <Layout title="Home">
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          fontSize: '20px',
          textAlign: 'center',
          height: "70vh",
          padding: "2em",
        }}>
	        <div style={{
				display: 'flex',
				maxWidth: "35em",
				justifyContent: 'center',
				alignItems: 'center',
				flexDirection: 'column',
	        }}>
				<h1 style={{
					fontSize: "50px",
				}}>
					Get your backend for free
		        </h1>
		        <p style={{
					color: "var(--ifm-color-primary-lighter)",
					marginBottom: "1.5em",
		        }}>
			        Aft gives you a modern API, scriptable RPCs, login, access controls and more. 
			        No need to run a database, or pay a penny.
		        </p>
	            <Link
	              className={clsx(
	                'button button--outline button--secondary button--lg'
	              )}
	              to={useBaseUrl('docs/')}>
	              Learn More
	            </Link>
	      </div>
      </div>
    </Layout>
  );
}

export default Hello;
