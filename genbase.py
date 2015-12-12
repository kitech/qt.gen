# encoding: utf8


class GenerateBase:
    def generateHeader(self, module):
        return

    def generateClasses(self, module, class_decls):
        return

    def generateClass(self, class_name, cs, methods):
        return

    def generateMethod(self, class_name, method_name, method_cursor):
        return

    # @return []
    def generateParams(self, class_name, method_name, method_cursor):
        return
    pass


# build the generated .swigcxx file
class TestBuilder:
    def tryBuild(self):
        return



class WarkerOne:
    def __init__(self): return

