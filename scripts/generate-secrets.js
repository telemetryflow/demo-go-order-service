#!/usr/bin/env node

/**
 * Generate Secure Secrets for order-service
 *
 * Usage:
 *   node scripts/generate-secrets.js              # Generate all secrets
 *   node scripts/generate-secrets.js --api-keys   # Generate only API keys
 *   node scripts/generate-secrets.js --jwt        # Generate only JWT secrets
 *   node scripts/generate-secrets.js --length 64  # Custom length
 *   node scripts/generate-secrets.js --format hex # Custom format
 *   node scripts/generate-secrets.js --env        # Output as .env format only
 */

const crypto = require('crypto');

// Parse arguments
const args = process.argv.slice(2);
let length = 32;
let format = 'base64';
let generateApiKeys = true;
let generateJwtSecrets = true;
let envOnly = false;

for (let i = 0; i < args.length; i++) {
  switch (args[i]) {
    case '--length':
      if (args[i + 1]) {
        length = parseInt(args[i + 1]);
        i++;
      }
      break;
    case '--format':
      if (args[i + 1]) {
        format = args[i + 1];
        i++;
      }
      break;
    case '--api-keys':
      generateApiKeys = true;
      generateJwtSecrets = false;
      break;
    case '--jwt':
      generateApiKeys = false;
      generateJwtSecrets = true;
      break;
    case '--env':
      envOnly = true;
      break;
    case '--help':
    case '-h':
      console.log(`
order-service - Secure Secret Generator

Usage:
  node scripts/generate-secrets.js [options]

Options:
  --length <number>   Length in bytes (default: 32)
  --format <format>   Output format: base64, hex, base64url (default: base64)
  --api-keys          Generate only TelemetryFlow API keys
  --jwt               Generate only JWT/Session secrets
  --env               Output in .env format only (no decorations)
  --help, -h          Show this help

Examples:
  node scripts/generate-secrets.js                    # Generate all secrets
  node scripts/generate-secrets.js --api-keys        # API keys only
  node scripts/generate-secrets.js --jwt --length 64 # JWT with 64 bytes
  node scripts/generate-secrets.js --env             # .env format output
`);
      process.exit(0);
  }
}

// Validation
if (length < 32) {
  console.error('Error: Length must be at least 32 bytes');
  process.exit(1);
}

const validFormats = ['base64', 'hex', 'base64url'];
if (!validFormats.includes(format)) {
  console.error(`Error: Format must be one of: ${validFormats.join(', ')}`);
  process.exit(1);
}

// Secret generation functions
function generateSecret(bytes, encoding) {
  const buffer = crypto.randomBytes(bytes);
  if (encoding === 'base64url') {
    return buffer.toString('base64')
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=/g, '');
  }
  return buffer.toString(encoding);
}

function generateHexString(length) {
  return crypto.randomBytes(Math.ceil(length / 2)).toString('hex').slice(0, length);
}

function generateApiKeyId() {
  return `tfk_${generateHexString(32)}`;
}

function generateApiKeySecret() {
  return `tfs_${generateHexString(64)}`;
}

// Generate secrets
const secrets = {};

if (generateApiKeys) {
  secrets.TELEMETRYFLOW_API_KEY_ID = generateApiKeyId();
  secrets.TELEMETRYFLOW_API_KEY_SECRET = generateApiKeySecret();
}

if (generateJwtSecrets) {
  secrets.JWT_SECRET = generateSecret(length, format);
  secrets.JWT_REFRESH_SECRET = generateSecret(length, format);
  secrets.SESSION_SECRET = generateSecret(length, format);
}

// Output
if (envOnly) {
  // Simple .env format output
  if (generateApiKeys) {
    console.log('# TelemetryFlow API Keys');
    console.log(`TELEMETRYFLOW_API_KEY_ID=${secrets.TELEMETRYFLOW_API_KEY_ID}`);
    console.log(`TELEMETRYFLOW_API_KEY_SECRET=${secrets.TELEMETRYFLOW_API_KEY_SECRET}`);
  }
  if (generateJwtSecrets) {
    if (generateApiKeys) console.log('');
    console.log('# JWT & Session Secrets');
    console.log(`JWT_SECRET=${secrets.JWT_SECRET}`);
    console.log(`JWT_REFRESH_SECRET=${secrets.JWT_REFRESH_SECRET}`);
    console.log(`JWT_EXPIRATION=24h`);
    console.log(`JWT_REFRESH_EXPIRATION=168h`);
    console.log(`SESSION_SECRET=${secrets.SESSION_SECRET}`);
  }
} else {
  // Decorated output
  console.log('\n=============================================');
  console.log('  order-service - Secret Generator');
  console.log('=============================================\n');

  if (generateApiKeys) {
    console.log('TelemetryFlow API Keys:');
    console.log('-----------------------');
    console.log(`  API Key ID:     ${secrets.TELEMETRYFLOW_API_KEY_ID}`);
    console.log(`  API Key Secret: ${secrets.TELEMETRYFLOW_API_KEY_SECRET}`);
    console.log('');
  }

  if (generateJwtSecrets) {
    console.log(`JWT & Session Secrets (${length} bytes, ${format}):`)
    console.log('----------------------------------------------');
    console.log(`  JWT Secret:         ${secrets.JWT_SECRET}`);
    console.log(`  JWT Refresh Secret: ${secrets.JWT_REFRESH_SECRET}`);
    console.log(`  Session Secret:     ${secrets.SESSION_SECRET}`);
    console.log('');
  }

  console.log('.env Format:');
  console.log('------------');
  if (generateApiKeys) {
    console.log(`TELEMETRYFLOW_API_KEY_ID=${secrets.TELEMETRYFLOW_API_KEY_ID}`);
    console.log(`TELEMETRYFLOW_API_KEY_SECRET=${secrets.TELEMETRYFLOW_API_KEY_SECRET}`);
  }
  if (generateJwtSecrets) {
    console.log(`JWT_SECRET=${secrets.JWT_SECRET}`);
    console.log(`JWT_REFRESH_SECRET=${secrets.JWT_REFRESH_SECRET}`);
    console.log(`JWT_EXPIRATION=24h`);
    console.log(`JWT_REFRESH_EXPIRATION=168h`);
    console.log(`SESSION_SECRET=${secrets.SESSION_SECRET}`);
  }
  console.log('');

  console.log('Docker Example:');
  console.log('---------------');
  const dockerEnvs = [];
  if (generateApiKeys) {
    dockerEnvs.push(`  -e TELEMETRYFLOW_API_KEY_ID="${secrets.TELEMETRYFLOW_API_KEY_ID}"`);
    dockerEnvs.push(`  -e TELEMETRYFLOW_API_KEY_SECRET="${secrets.TELEMETRYFLOW_API_KEY_SECRET}"`);
  }
  if (generateJwtSecrets) {
    dockerEnvs.push(`  -e JWT_SECRET="${secrets.JWT_SECRET}"`);
    dockerEnvs.push(`  -e SESSION_SECRET="${secrets.SESSION_SECRET}"`);
  }
  console.log(`docker run -d \\\n${dockerEnvs.join(' \\\n')} \\`);
  console.log('  order-service:latest\n');

  console.log('Security Tips:');
  console.log('--------------');
  console.log('- Never commit secrets to git');
  console.log('- Use different secrets per environment');
  console.log('- Rotate secrets every 90 days');
  console.log('- Store in secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.)');
  console.log('');
}

process.exit(0);
