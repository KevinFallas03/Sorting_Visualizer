def sort(array=[29, 20, 27, 14, 8, 25, 5, 20, 18, 1, 14, 26, 28, 6, 25, 1, 7, 11, 3, 7]):
    """Sort the array by using quicksort."""

    less = []
    equal = []
    greater = []

    print(array)
    
    if len(array) > 1:
        pivot = array[0]
        for x in array:
            if x < pivot:
                less.append(x)
            elif x == pivot:
                equal.append(x)
            elif x > pivot:
                greater.append(x)
        # Don't forget to return something!
        return (sort(less)+equal+sort(greater))  # Just use the + operator to join lists
    # Note that you want equal ^^^^^ not pivot
    else:  # You need to handle the part at the end of the recursion - when you only have one element in your array, just return the array.
        return (array)
    
    
lista = sort([29, 20, 27, 14, 8, 25, 5, 20, 18, 1, 14, 26, 28, 6, 25, 1, 7, 11, 3, 7])
print(lista)