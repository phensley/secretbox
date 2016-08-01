
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

## Acknowledgements

Derived from [levigross/keylesscrypto](https://github.com/levigross/keylesscrypto)

