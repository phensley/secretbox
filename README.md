
# secretbox

Encrypt / decrypt files using Shamir secret sharing.


## Encrypt

    $ echo "hello world" >foo.txt
    $ secretbox encrypt -p 4 -t 2 -i foo.txt -o foo.crypted -e base64

	Generated secret key in 4 parts with thresdhold 2:

	 [1] GfrtqFq8sYWCWggF5cO8EVjFAzpIA18p7CeqC9IyovCgp+yBW1v61ikQMGpyuwhVZTKVdNlEgsAB
	 [2] uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC
	 [3] Li9jwAunsExtEJtCnYB/VxbdPjGiKPmYEVziW89r5F1kINtQWYZr6HX/mlxvzEMazNGkWeRMRYQD
	 [4] 4aytTF0KPuW96/ElKenNvuP56qTyzVuF/mEeg2shDWpRfBTdXhkGtb8vKh3L4zZh9RNitTBQ+WoE

	Encrypting to 'foo.crypted'

	Success!


## Decrypt

	$ secretbox decrypt -i foo.crypted -o foo.decrypted -e base64

Interactively enter the Shamir secrets and decrypt the file:

	This interactive shell will allow you to enter the key parts.

	Commands:

	 add <part>    - adds a key part
	 list          - view the parts that have been entered
	 del <num>     - deletes the part by #
	 done          - indicate you have entered parts and are ready to decrypt

	 exit          - exit immediately without decrypting
	 help          - display this message

	>>

Add at least `threshold` parts:

	>> add uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC
	add: uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC

	>> add 4aytTF0KPuW96/ElKenNvuP56qTyzVuF/mEeg2shDWpRfBTdXhkGtb8vKh3L4zZh9RNitTBQ+WoE
	add: 4aytTF0KPuW96/ElKenNvuP56qTyzVuF/mEeg2shDWpRfBTdXhkGtb8vKh3L4zZh9RNitTBQ+WoE

Use up-arrow / `del` to edit / re-enter a part:

	>> list
	 [1] uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC
	 [2] 4aytTF0KPuW96/ElKenNvuP56qTyzVuF/mEeg2shDWpRfBTdXhkGtb8vKh3L4zZh9RNitTBQ+WoE

	>> del 2
	deleting key at index 2

	>> list
	 [1] uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC

	>> add Li9jwAunsExtEJtCnYB/VxbdPjGiKPmYEVziW89r5F1kINtQWYZr6HX/mlxvzEMazNGkWeRMRYQD
	add: Li9jwAunsExtEJtCnYB/VxbdPjGiKPmYEVziW89r5F1kINtQWYZr6HX/mlxvzEMazNGkWeRMRYQD

	>> list
	 [1] uMgk9K4nPaWXNV/soSyTdDHRrbnXsKpN4uzGc0zKx4YG7k21WGWu91sFz0fseuuwFS0xwndIq6YC
	 [2] Li9jwAunsExtEJtCnYB/VxbdPjGiKPmYEVziW89r5F1kINtQWYZr6HX/mlxvzEMazNGkWeRMRYQD

When ready, type `done` to decrypt:

	>> done
	2 keys entered. validating.

	Success!

	$ cat foo.decrypted
	hello world

## Encrypting secrets using a PGP keyring

Create an ASCII-armored keyring containing the PGP keys for the recipients of the secrets.
If you have the public keys but don't have them on a keyring, you can create a temporary
keyring:

    $ gpg --no-default-keyring --keyring ./recipients.gpg --import bob-pub.asc alice-pub.asc bill-pub.asc

Then export the list of recipients you want to receive the secrets:

    $ gpg --no-default-keyring --keyring ./recipients.gpg --export -a \
		bob@example.com alice@example.com bill@example.com >keyring.asc

Finally, use `secretbox` to encrypt a file, and then encrypt the secrets for the recipients.

Note: The number of secret parts must equal the number of keys in the keyring.  In this case we're
splitting the secret into 3 parts and have 3 recipients.

    $ secretbox encrypt -p 3 -t 2 -i foo.txt -o foo.crypted -k keyring.asc >secrets.txt
	$ cat secrets.txt

	Encrypting to 'foo.crypted'
	Encrypted using secret key in 3 parts with threshold 2:

	Encrypting secret 1 for recipient:
	 --> Bob Jones <bob@example.com>

	-----BEGIN PGP MESSAGE-----
	...
	-----END PGP MESSAGE-----

	Encrypting secret 2 for recipient:
	 --> Alice Smith <alice@example.com>

	-----BEGIN PGP MESSAGE-----
	...
	-----END PGP MESSAGE-----

	Encrypting secret 3 for recipient:
	 --> Bill Davis <bill@example.com>

	-----BEGIN PGP MESSAGE-----
	...
	-----END PGP MESSAGE-----

Now you can snip out each encrypted message and send it to the recipients.


## Acknowledgements

Derived from [levigross/keylesscrypto](https://github.com/levigross/keylesscrypto)

