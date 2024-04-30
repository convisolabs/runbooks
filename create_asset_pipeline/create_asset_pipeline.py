  - apt-get update && apt-get install -y jq

create_asset:
  script:
    - |
      PROJECT_NAME=$(basename "$CI_PROJECT_DIR")
      echo "Project name: $PROJECT_NAME"
      RESPONSE=$(curl -s -X POST \
        -H "x-api-key: $CONVISO_API_KEY" \
        -H "Content-Type: application/json" \
        -d '{"query": "{ assets(filter: { name: { eq: \"'$PROJECT_NAME'\" } }) { id } }"}' \
        https://app.convisoappsec.com/graphql)
      ASSET_ID=$(echo "$RESPONSE" | jq -r '.data.assets[0].id')
      if [ -n "$ASSET_ID" ] && [ "$ASSET_ID" != "null" ]; then
        echo "Asset already exists with ID: $ASSET_ID"
      else
        echo "Creating asset..."
        RESPONSE=$(curl -s -X POST \
          -H "x-api-key: $CONVISO_API_KEY" \
          -H "Content-Type: application/json" \
          -d '{"query": "mutation { createAsset(input: { companyId: 600, name: \"'$PROJECT_NAME'\", businessImpact: HIGH, dataClassification: NON_SENSITIVE, description: \"arquivo demo33333\" }) { asset { id name businessImpact description } clientMutationId errors } }"}' \
          https://app.convisoappsec.com/graphql)
        ASSET_ID=$(echo "$RESPONSE" | jq -r '.data.createAsset.asset.id')
        if [ -n "$ASSET_ID" ] && [ "$ASSET_ID" != "null" ]; then
          echo "Asset created successfully with ID: $ASSET_ID"
        else
          echo "Failed to create asset. Response: $RESPONSE"
        fi
      fi

#conviso-ast:
#  image: convisoappsec/convisocli:latest
#  services:
#    - docker:dind
#  only:
#    variables:
#      - $CONVISO_API_KEY
#  script:
#    - conviso ast run
#  tags:
#    - docker