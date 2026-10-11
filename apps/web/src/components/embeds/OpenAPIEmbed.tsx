import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';
import type { EmbedProps } from './types.ts';

/** Renders an OpenAPI 3.x or Swagger 2.0 document as API reference. Calls to the API are disabled. */
export default function OpenAPIEmbed({ src }: EmbedProps) {
  return (
    <div className="rounded-lg border border-line">
      <SwaggerUI
        url={src}
        supportedSubmitMethods={[]}
        docExpansion="list"
        defaultModelsExpandDepth={0}
        tryItOutEnabled={false}
      />
    </div>
  );
}
