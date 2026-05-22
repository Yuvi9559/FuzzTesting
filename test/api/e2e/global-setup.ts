import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log('🔧 Setting up fuzztesting e2e test environment...');
  
  // Launch a browser to verify we can connect
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  try {
    // Wait for fuzztesting API to be available
    const baseURL = process.env.fuzztesting_API_URL || 'http://localhost:8080/api/v1';
    const maxRetries = 30;
    let retries = 0;
    
    while (retries < maxRetries) {
      try {
        const response = await page.request.get(`${baseURL}/health`);
        if (response.ok()) {
          console.log('✅ fuzztesting API is ready');
          break;
        }
      } catch (error) {
        // Continue retrying
      }
      
      retries++;
      if (retries >= maxRetries) {
        throw new Error(`fuzztesting API not available at ${baseURL}/health after ${maxRetries} retries`);
      }
      
      console.log(`⏳ Waiting for fuzztesting API (attempt ${retries}/${maxRetries})...`);
      await new Promise(resolve => setTimeout(resolve, 2000));
    }
    
    // Verify readiness endpoint
    try {
      const readinessResponse = await page.request.get(`${baseURL}/readiness`);
      if (readinessResponse.ok()) {
        console.log('✅ fuzztesting system is ready');
      } else {
        console.warn('⚠️  fuzztesting system readiness check failed, proceeding anyway');
      }
    } catch (error) {
      console.warn('⚠️  Could not check system readiness, proceeding anyway');
    }
    
  } finally {
    await browser.close();
  }
  
  console.log('✅ Global setup completed');
}

export default globalSetup;