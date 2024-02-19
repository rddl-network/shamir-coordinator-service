# Import the required module. This might need to be installed first.
# You can install it using pip if it's available.
import sys
import shamir_mnemonic

arguments = sys.argv[1:]
# Call the combine_mnemonics function and convert the result to hex

try:
    result = shamir_mnemonic.combine_mnemonics(arguments).hex()
    print(result)
except Exception as e:
    print(f"An error occurred: {str(e)}")
    exit(1)

