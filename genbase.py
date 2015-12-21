# encoding: utf8

from genutil import *

class GenerateBase(object):
    def __init__(self):
        super(GenerateBase, self).__init__()
        self.gctx = None
        self.gutil = GenUtil()
        return

    def setGenContext(self, ctx):
        self.gctx = ctx
        return

    def genpass(self, module):
        raise 'not impled'
        return

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

