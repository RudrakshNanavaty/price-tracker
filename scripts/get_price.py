try:
	from sys import argv
	from selenium import webdriver
	from selenium.webdriver.common.by import By
	import logging
except ImportError:
	logging.error("Dependencies missing. Please run 'pip install -r ./scripts/requirements.txt'")
	exit(1)

# set up logging in ./scripts/get_price.log
logging.basicConfig(filename='./scripts/get_price.log', level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

try:
	# initialise Firefox webdriver with options
	options = webdriver.FirefoxOptions()
	# don't open a window
	options.add_argument("-headless")
	driver = webdriver.Firefox(options=options)
except:
	logging.error("Error initialising Firefox webdriver.")
	exit(1)

try:
	driver.get(argv[1])
except:
	logging.error("Error fetching page from URL.")
	exit(1)

try:
	# get the price of the product
	price_string = driver.find_element(By.CLASS_NAME, 'a-price').text
except:
	logging.error("Error fetching price from URL.")
	exit(1)

driver.quit()

# print the price without the currency symbol or commas
print(float(price_string.replace('₹', '').replace(',', '').replace('\n', '.')), end="")